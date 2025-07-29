import sys
import os
from cx_Freeze import setup, Executable

# Пути к дополнительным ресурсам из clip
clip_pkg_path = r"C:\Users\Sssemenogka\AppData\Local\Programs\Python\Python312\Lib\site-packages\clip"

# Включаем необходимые файлы из clip
include_files = [
    (os.path.join(clip_pkg_path, "bpe_simple_vocab_16e6.txt.gz"), os.path.join("clip", "bpe_simple_vocab_16e6.txt.gz")),
    ("alldata/ndata", "alldata/ndata"),  # если нужны ваши json-файлы
]

# Опции сборки
build_exe_options = {
    "packages": ["torch", "clip", "PIL", "requests", "PyQt6", "numpy"],
    "include_files": include_files,
    "excludes": ["tkinter"],  # исключаем, если не нужен
}

# GUI-базовый режим для Windows (чтобы не было консоли)
base = None
if sys.platform == "win32":
    base = "Win32GUI"

setup(
    name="LampSearch",
    version="1.0",
    description="Приложение для поиска светильников",
    options={"build_exe": build_exe_options},
    executables=[Executable("app.py", base=base)]
)
