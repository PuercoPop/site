use axum::{extract::Multipart, response::Html};

pub async fn index() -> Html<&'static str> {
    Html("こんにちは<a href=/upload>Upload</a>")
}

pub async fn form() -> Html<&'static str> {
    Html("<html><body><h1>Upload Image</h1><form action=/upload method=post enctype='multipart/form-data'>
<input name=media type=file />
<input type=submit value='upload'/>
</form></body></html>")
}

pub async fn upload(mut multipart: Multipart) {
    // TODO: Replace unwrap calls with ?
    while let Some(field) = multipart.next_field().await.unwrap() {
        let name = field.name().unwrap().to_string();
        println!("Field is {}", name);
        println!("Filename is {}", field.file_name().unwrap().to_string());
        // println!("Content Type is {}", field.content_type()?.to_string());
    }
}
