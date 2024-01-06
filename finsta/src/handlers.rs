use crate::HTTPContext;
use axum::{
    extract::{Multipart, State},
    response::{Html, IntoResponse, Redirect, Response, Result as HandlerResult},
    Form,
};
use minijinja::context;
use serde::Deserialize;
use std::sync::Arc;
use tokio_postgres::Error as PgError;

#[derive(thiserror::Error, Debug)]
pub enum HandlerError {
    #[error(transparent)]
    TemplateError(#[from] minijinja::Error),
    #[error(transparent)]
    DBError(#[from] PgError),
    #[error(transparent)]
    HTTPError(#[from] http::Error),
}

impl IntoResponse for HandlerError {
    fn into_response(self) -> Response {
        (http::StatusCode::INTERNAL_SERVER_ERROR, self.to_string()).into_response()
    }
}

pub async fn index(
    State(state): State<Arc<HTTPContext>>,
) -> HandlerResult<Html<String>, HandlerError> {
    let tmpl = state.templates.get_template("index.html")?;
    Ok(Html(tmpl.render(context!())?))
}

#[derive(Deserialize)]
pub(crate) struct SignInQP {
    email: String,
    password: String,
}

// TODO: Track sign-in attempts
pub(crate) async fn sign_in(
    State(state): State<Arc<HTTPContext>>,
    Form(params): Form<SignInQP>,
) -> HandlerResult<impl IntoResponse, HandlerError> {
    // 1. ✔ Extract request parameters
    // 2. ✔ Check against the database
    // 3. If successful redirect
    let is_valid =
        crate::users::authenticate_user(&state.db, params.email, params.password).await?;

    if is_valid {
        // TODO: Get URL from query param, redirect_to
        // https://github.com/tokio-rs/axum/discussions/351
        Ok(Redirect::temporary("/"))
    } else {
        // TODO: Figure out how to return a response on this branch. Seems I
        // have to return the same type on both branches?
        // let response = Response::builder().status(StatusCode::NOT_FOUND).body(
        //     "This are not the droids that you are looking for",
        // )?;
        // Ok(response)
        Ok(Redirect::temporary("/"))
    }
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
