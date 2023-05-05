//! The template API
use minijinja::{Environment};

pub enum Error {

}

/// Initializes the minijinja environment and loads all the templates. Takes the
/// directory where the templates reside as a path.
pub fn new() -> Result<Environment, Error> {
    let env = Environment::new();
    env.add_template("layout",include_str!("../templates/layout.html"))?;
    Ok(env)
}


/// Renders a template
pub fn render(name: String) -> Result<String, Error> {
    todo!()
}
