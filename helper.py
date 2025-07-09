import os
import sys
import hashlib
import shutil

def delete():
    file_path = resource_path("cache.pt")
    dir_path = resource_path("alldata")
    if os.path.isfile(file_path):
      os.remove(file_path)
      print(f"Файл {file_path} удалён")
    if os.path.isdir(dir_path):
      shutil.rmtree(dir_path)
      print(f"Папка {dir_path} и всё её содержимое удалены")


def file_sha256(path):
    if not os.path.exists(path):
        return None
    with open(path, "rb") as f:
        return hashlib.sha256(f.read()).hexdigest()
    
def resource_path(relative_path):
    if hasattr(sys, '_MEIPASS'):
        return os.path.join(sys._MEIPASS, relative_path)
    return os.path.join(os.path.abspath("."), relative_path)

def is_number(s):
    try:
        float(s)
        return True
    except ValueError:
        return False