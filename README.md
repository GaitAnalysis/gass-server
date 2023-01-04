# gass-server

## Requirements
TODO: Automate in setup

**Python 3.10.2 [torch not available yet for 3.11 and beyond]**

```console
pip install -r requirements.txt

wget -P AlphaPose/detector/yolo/data https://www.dropbox.com/s/udbe7w7kp0un7ob/yolov3-spp.weights

mkdir AlphaPose/detector/tracker/data

wget -P AlphaPose/detector/tracker/data https://www.dropbox.com/s/hnyoa8q5ld5vanj/JDE-1088x608-uncertainty

wget -P AlphaPose/pretrained_models https://www.dropbox.com/s/f5ai60e7zd7i4h9/fast_res50_256x192.pth

wget -P AlphaPose/detector/yolox/data/ https://github.com/Megvii-BaseDetection/YOLOX/releases/download/0.1.0/yolox_x.pth
```

## API Docs
http://0.0.0.0:8000/docs

## Run
```console
$ python main.py
```