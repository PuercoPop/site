use axum::{
    extract::Multipart,
    response::Html,
    routing::{get, post},
    Router,
};

pub fn new() -> Router {
    Router::new()
        .route("/", get(index))
        .route("/upload", get(form))
        .route("/upload", post(upload))
}

pub async fn index() -> Html<&'static str> {
    Html("こんにちは<a href=/upload>Upload</a>")
}

// TODO: whoami page/endpoint

pub async fn form() -> Html<&'static str> {
    Html("<html><body><h1>Upload Image</h1><form action=/upload method=post enctype='multipart/form-data'>
<input name=media type=file />
<input type=submit value='upload'/>
</form></body></html>")
}

pub async fn upload(mut multipart: Multipart) {
    // TODO: Replace unwrap calls with ?
    // Handler should return something like Result<(), axum::http::StatusCode>
    while let Some(field) = multipart.next_field().await.unwrap() {
        let name = field.name().unwrap().to_string();
        println!("Field is {}", name);
        println!("Filename is {}", field.file_name().unwrap().to_string());
        // println!("Content Type is {}", field.content_type()?.to_string());
    }
}
