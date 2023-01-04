#!/bin/bash
echo 'Installing requirements...'
pip install -r requirements.txt

echo 'Getting AlphaPose...'
rm -rf ./AlphaPose
git clone https://github.com/MVIG-SJTU/AlphaPose.git

echo 'Getting parameters...'
mkdir AlphaPose/detector/yolo/data
wget -P AlphaPose/detector/yolo/data https://www.dropbox.com/s/udbe7w7kp0un7ob/yolov3-spp.weights

mkdir AlphaPose/detector/tracker/data
wget -P AlphaPose/detector/tracker/data https://www.dropbox.com/s/hnyoa8q5ld5vanj/JDE-1088x608-uncertainty

wget -P AlphaPose/pretrained_models https://www.dropbox.com/s/f5ai60e7zd7i4h9/fast_res50_256x192.pth

wget -P AlphaPose/detector/yolox/data/ https://github.com/Megvii-BaseDetection/YOLOX/releases/download/0.1.0/yolox_x.pth

#echo 'Building AlphaPose...'
#cd AlphaPose
#python setup.py build develop
#cd ..
