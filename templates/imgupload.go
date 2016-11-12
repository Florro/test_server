package templates


var ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body><img src="data:image/jpg;base64,{{.Image}}"></body>`

var UploadTemplate string = `<html>
<head>
    <title>Upload file</title>
<form enctype="multipart/form-data" action="http://127.0.0.1:8080/tf2" method="post">
      <input type="file" name="uploadfile" />
      <input type="hidden" name="token" value="{{.}}"/>
      <input type="submit" value="upload" />
</head>
<body>
</form>
</body>
</html>`


var UploadTemplate2 string = `<html>
<head>
    <title>Upload file</title>
<form enctype="multipart/form-data" action="http://127.0.0.1:8080/all" method="post">
      <input type="file" name="uploadfile" />
      <input type="hidden" name="token" value="{{.}}"/>
      <input type="submit" value="upload" />
</head>
<body>
</form>
</body>
</html>`
