use crate::db::{all_posts, post_by_slug, posts_by_tag, recent_posts, tags_count};
use crate::Context;
use axum::{
    extract::{Path as URLPath, State},
    http,
    response::{Html, IntoResponse, Response, Result as HandlerResult},
};
use http::{header, HeaderMap};
use minijinja::context;
use std::sync::Arc;

/// HandlerError ocurr within HTTP handlers
#[derive(Debug, thiserror::Error)]
pub enum HandlerError {
    #[error(transparent)]
    TemplateError(#[from] minijinja::Error),
    #[error(transparent)]
    DBError(#[from] tokio_postgres::Error),
    #[error(transparent)]
    InvalidHeader(#[from] hyper::header::InvalidHeaderValue),
}

impl IntoResponse for HandlerError {
    fn into_response(self) -> Response {
        (http::StatusCode::INTERNAL_SERVER_ERROR, self.to_string()).into_response()
    }
}

#[axum::debug_handler]
pub(crate) async fn index(
    State(state): State<Arc<Context>>,
) -> HandlerResult<Html<String>, HandlerError> {
    let tmpl = state.templates.get_template("index.html")?;
    let posts = recent_posts(&state.db).await.expect("IOU a ?");
    Ok(Html(tmpl.render(context!(latest_posts => posts))?))
}

#[axum::debug_handler]
pub(crate) async fn show_post(
    State(state): State<Arc<Context>>,
    URLPath(slug): URLPath<String>,
) -> HandlerResult<Html<String>, HandlerError> {
    let tmpl = state.templates.get_template("post.html")?;

    let post = post_by_slug(&state.db, slug).await?;
    Ok(Html(tmpl.render(context!(post => post))?))
}

#[axum::debug_handler]
pub(crate) async fn list_tags(
    State(state): State<Arc<Context>>,
) -> HandlerResult<Html<String>, HandlerError> {
    let tmpl = state.templates.get_template("tag-list.html")?;
    let tags = tags_count(&state.db).await?;
    Ok(Html(tmpl.render(context!(tags => tags))?))
}

#[axum::debug_handler]
pub(crate) async fn show_tag(
    State(state): State<Arc<Context>>,
    URLPath(tag): URLPath<String>,
) -> HandlerResult<Html<String>, HandlerError> {
    let tmpl = state.templates.get_template("tag-detail.html")?;
    let posts = posts_by_tag(&state.db, &tag).await?;
    Ok(Html(tmpl.render(context!(tag => tag, posts => posts))?))
}

// TODO: Add updated as a top-level tag to the feed.
// TODO: After we stop recreating the database on each deploy, implement a
// paginated feed.
// TODO: After we feed is paginated we can include the posts content in the
// feed.
/// Implements the blog's Atom feed. See:
/// https://datatracker.ietf.org/doc/html/rfc4287
/// and https://github.com/rackerlabs/riss/blob/master/cookbook/atom-feed-paging-and-archiving.md
#[axum::debug_handler]
pub(crate) async fn feed(
    State(state): State<Arc<Context>>,
) -> HandlerResult<Response, HandlerError> {
    let tmpl = state.templates.get_template("atom.xml")?;
    let posts = all_posts(&state.db).await?;
    let mut headers = HeaderMap::new();
    headers.insert(header::CONTENT_TYPE, "application/atom+xml".parse()?);
    let body = tmpl.render(context!(posts => posts))?;
    Ok((headers, body).into_response())
}
