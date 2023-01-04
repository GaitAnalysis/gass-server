from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()

origins = [
    'http://localhost:3000',
    'localhost:3000'
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=['*'],
    allow_headers=['*']
)


@app.get('/')
async def root():
    return {'message': 'success'}


@app.post('/upload/video/')
async def upload_video(file: UploadFile):
    print('IN')
    print(file.content_type)
    if not file.content_type.startswith('video'):
        raise HTTPException(status_code=400, detail='Wrong file format, must be video.')
    content = await file.read()
    filename = file.filename
    with open(filename, 'wb') as f:
        f.write(content)
    return {'filename': filename}
