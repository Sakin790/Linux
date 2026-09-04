const express = require('express');
const http = require('http');
const WebSocket = require('ws');
const app = express();
const PORT = 5000;

// ইনকাস ইউনিক্স সকেটের পাথ
const INCUS_SOCKET = '/var/lib/incus/unix.socket';

app.use(express.json());
app.use(express.static('public'));

// ১. কন্টেইনারের লিস্ট ও রিয়েল-টাইম স্ট্যাটাস দেখার API
app.get('/api/containers', (req, res) => {
    const options = {
        socketPath: INCUS_SOCKET,
        path: '/1.0/instances?recursion=1',
        method: 'GET'
    };

    const clientRequest = http.request(options, (authRes) => {
        let data = '';
        authRes.on('data', (chunk) => data += chunk);
        authRes.on('end', () => {
            try {
                const responseJSON = JSON.parse(data);
                res.json({ success: true, containers: responseJSON.metadata || [] });
            } catch (error) {
                res.status(500).json({ success: false, error: "Parsing error" });
            }
        });
    });

    clientRequest.on('error', (e) => res.status(500).json({ success: false, error: e.message }));
    clientRequest.end();
});

// ২. কন্টেইনার স্টার্ট বা স্টপ করার API
app.put('/api/containers/:name/state', (req, res) => {
    const containerName = req.params.name;
    const { action } = req.body;

    const postData = JSON.stringify({ action: action, timeout: 30 });

    const options = {
        socketPath: INCUS_SOCKET,
        path: `/1.0/instances/${containerName}/state`,
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Content-Length': Buffer.byteLength(postData)
        }
    };

    const clientRequest = http.request(options, (authRes) => {
        let data = '';
        authRes.on('data', (chunk) => data += chunk);
        authRes.on('end', () => res.json({ success: true, info: JSON.parse(data) }));
    });

    clientRequest.on('error', (e) => res.status(500).json({ success: false, error: e.message }));
    clientRequest.write(postData);
    clientRequest.end();
});

const server = http.createServer(app);
const wss = new WebSocket.Server({ noServer: true });

// HTTP সার্ভারে আপগ্রেড রিকোয়েস্ট হ্যান্ডেল করা
server.on('upgrade', (request, socket, head) => {
    const pathname = request.url;

    if (pathname.startsWith('/terminal/')) {
        wss.handleUpgrade(request, socket, head, (ws) => {
            wss.emit('connection', ws, request);
        });
    } else {
        socket.destroy();
    }
});

wss.on('connection', (ws, req) => {
    const parts = req.url.split('/');
    const containerName = parts[parts.length - 1];

    if (!containerName) {
        ws.close(1008, "Container name is required");
        return;
    }

    // bash এবং ইন্টারেক্টিভ মোড (-i) ফোর্স করা হয়েছে যাতে শেলটি হুট করে বন্ধ না হয়
    const postData = JSON.stringify({
        command: ["/bin/bash", "-i"],
        environment: { "HOME": "/root", "TERM": "xterm-256color" },
        "interactive": true,
        "wait-for-websocket": true
    });

    const options = {
        socketPath: INCUS_SOCKET,
        path: `/1.0/instances/${containerName}/exec`,
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Content-Length': Buffer.byteLength(postData)
        }
    };

    const clientRequest = http.request(options, (res) => {
        let responseData = '';
        res.on('data', (chunk) => responseData += chunk);
        res.on('end', () => {
            try {
                const result = JSON.parse(responseData);
                if (result.metadata && result.metadata.metadata && result.metadata.metadata.fds) {
                    connectToIncusTerminal(result.metadata.id, result.metadata.metadata.fds['0'], ws);
                } else {
                    ws.send("\r\n[Error] Interactive shell initialization failed. Is container running?\r\n");
                    ws.close();
                }
            } catch (err) {
                ws.close();
            }
        });
    });

    clientRequest.on('error', () => ws.close());
    clientRequest.write(postData);
    clientRequest.end();
});

function connectToIncusTerminal(opId, secret, browserWs) {
    const incusWs = new WebSocket(
        `ws+unix://${INCUS_SOCKET}:/1.0/operations/${opId}/websocket?secret=${secret}`
    );

    // আইডল কানেকশন ড্রপ রোধে হার্টবিট সেটআপ
    const heartbeat = setInterval(() => {
        if (browserWs.readyState === WebSocket.OPEN) {
            browserWs.ping();
        }
    }, 25000); // প্রতি ২৫ সেকেন্ডে পিং পাঠানো হবে

    incusWs.on('open', () => {
        browserWs.on('message', (message) => {
            if (incusWs.readyState === WebSocket.OPEN) {
                incusWs.send(message);
            }
        });
    });

    incusWs.on('message', (data) => {
        if (browserWs.readyState === WebSocket.OPEN) {
            browserWs.send(data);
        }
    });

    const cleanUp = () => {
        clearInterval(heartbeat);
        if (browserWs.readyState === WebSocket.OPEN) browserWs.close();
        if (incusWs.readyState === WebSocket.OPEN) incusWs.close();
    };

    incusWs.on('close', cleanUp);
    incusWs.on('error', (err) => {
        console.error('Incus terminal socket error:', err.message);
        if (browserWs.readyState === WebSocket.OPEN) {
            browserWs.send('\r\n[Error] Incus terminal socket error.\r\n');
        }
        cleanUp();
    });

    browserWs.on('close', cleanUp);
}

server.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});
