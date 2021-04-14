from PIL import Image
from io import BytesIO
import base64
import time
def handler(event):
    startTime = int(round(time.time() * 1000))
    img = Image.open(BytesIO(base64.b64decode(event['img'])))
    img = img.resize((event['height'], event['width']), Image.ANTIALIAS)
    buffer = BytesIO()
    img.save(buffer, 'JPEG')
    buffer.seek(0)
    
    return {'img': str(base64.b64encode(buffer.read()), encoding='ascii'), 
            'startTime': startTime, 
            'retTime': int(round(time.time() * 1000)), 
            'invokeTime': startTime}

if __name__ == "__main__":
    f = open('test.jpg', 'rb')
    # img = Image.open('test.jpg')
    # img = img.resize((200, 200), Image.ANTIALIAS)
    # t = open('t.jpg', 'w')
    # img.save(t, 'JPEG')
    print(handler({'img': str(base64.b64encode(f.read()), encoding='ascii'), 'height': 200, 'width': 200}))
