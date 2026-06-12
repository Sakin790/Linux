import subprocess
result=subprocess.run(["pwd"],capture_output=True,text=True)

if result.returncode==0:
    print("Command executed successfully")
    print("Output:",result.stdout)
else:  
     print("Command execution failed")
     print("Error:",result.stderr)
