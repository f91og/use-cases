use serde::{Serialize, Deserialize};

// 过程宏3种形式之一：派生宏(derive macros)- #[derive(CustomDerive)]
// 简而言之就是为结构体或者枚举生成特征的实现
// https://www.bilibili.com/video/BV1hp4y1k7SV?p=107&vd_source=6e0c18349c9d2d60a1a7c15986b114e6
#[derive(Debug, Serialize, Deserialize)]
struct User {
    name: String,
    age: i32,
    gender: String,
    friends: Vec<String> // 不定长数组，类似于C++中的Vector
}

fn main() {
    // read yaml
    let yaml_str = include_str!("./file.yaml");
    let mut user: User = serde_yaml::from_str(yaml_str).expect("yaml read failed!");

    println!("{:?}", user); // 打印结构体的方式，需要使用#[derive(Debug)]这个宏，{}是占位符，用{:?}来打印结构体，如果属性较多的话可以使用另一个占位符 {:#?}
    // User { name: "Alex", age: 21, gender: "male", friends: ["Bob", "Alice", "Tide"] }

    // update and write
    user.name = String::from("Bob");
    let f = std::fs::OpenOptions::new().write(true).create(true).open("file2.yaml").expect("Couldn't open file");
    serde_yaml::to_writer(f, &user).unwrap();
}

//offical demo: https://github.com/dtolnay/serde-yaml
//use serde::{Serialize, Deserialize};
//
//#[derive(Debug, PartialEq, Serialize, Deserialize)]
//struct Point {
//    x: f64,
//    y: f64,
//}
//
//fn main() -> Result<(), serde_yaml::Error> {
//    let point = Point { x: 1.0, y: 2.0 };
//
//    let yaml = serde_yaml::to_string(&point)?;
//    assert_eq!(yaml, "x: 1.0\ny: 2.0\n");
//
//    let deserialized_point: Point = serde_yaml::from_str(&yaml)?;
//    assert_eq!(point, deserialized_point);
//    Ok(())
//}
