import json
import re
import sys

p = re.compile(r'filename="([\w\W]+?)"')

def fileBody(fName):
  with open(fName, 'r') as file:
    return file.read()

def replace(data):
  for m in re.finditer(p, data):
    data = re.sub(m.group(), fileBody(m.group(1)), data)
  return data

def checkArgs():
  if len(sys.argv) < 2:
    print("Use: python", sys.argv[0], "mail.txt")
    quit(1)
  
  return sys.argv[1]

def main():
  fName = checkArgs()

  data = replace(fileBody(fName))
  j = {"Data": data}
  print(json.dumps(j))

main()
