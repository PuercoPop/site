#[derive(Debug)]
struct Post {
    title: String,
    // pubdate: chrono::NaiveDate,
    // draft: boolean,
    tags: Vec<Tag>
}

#[derive(Debug)]
struct Tag {
    name: String
}

fn main() {
    println!("Hello, world!");
}
