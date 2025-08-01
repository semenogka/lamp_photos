import torch
import os
import requests
import io
import gzip
_cache = {}



def load_cache(buttons_files):
    global _cache
    for category in buttons_files:
      filename = os.path.basename(buttons_files[category])
      raw_url = f"https://raw.githubusercontent.com/semenogka/lamp_photos/main/cacheGZ/{filename}.gz"
      print(f"Загрузка: {category}")
      res = requests.get(raw_url)
      if res.status_code == 200:
            try:
                  buffer = io.BytesIO(res.content)
                  with gzip.GzipFile(fileobj=buffer, mode="rb") as gz:
                        entry = torch.load(gz, map_location="cpu")
                  _cache[category] = entry
            except Exception as e:
                  print(f"[!] Ошибка torch.load({category}): {e}")
                  continue
      else:
           print(f"[!] Ошибка скачивания({category})")

    return _cache