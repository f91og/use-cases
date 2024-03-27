use std::collections::VecDeque;

fn main() {
    // Check if XML string is provided as command line argument
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        println!("Please provide an XML string as the first command line argument.");
        return;
    }
    if args.len() > 2 {
        println!("Too many command line arguments provided, only the first argument (XML string) is considered.");
        return;
    }

    let xml_string = &args[1];

    // Check if the XML string is valid
    if is_valid_xml(xml_string) {
        print!("Valid");
    } else {
        println!("Invalid");
    }
}

// Function to check if the XML string is valid
fn is_valid_xml(xml: &str) -> bool {
    // Check if the string contains '<' and '>', and their count matches
    if !xml.contains('<') || !xml.contains('>') {
        return false;
    }

    let mut stack = VecDeque::new();

    let mut i = 0;
    while i < xml.len() {
        if xml[i..].starts_with('<') {
            if let Some(end_tag_index) = xml[i + 1..].find('>') {
                let tag = &xml[i + 1..i + 1 + end_tag_index];

                // Check if it's an opening tag
                if !tag.starts_with('/') {
                    stack.push_back(tag.to_string());
                } else {
                    // Check if it's a closing tag
                    let opening_tag = stack.pop_back();
                    if opening_tag.is_none() || &opening_tag.unwrap() != &tag[1..] {
                        return false; // Mismatched closing tag
                    }
                }

                // i + 1 + end_tag_index + 1 points to the character after ">"
                i += 2 + end_tag_index;
            } else {
                // No matching '>' found
                return false;
            }
        } else {
            // Skip non-tag characters
            i += 1;
        }
    }

    stack.is_empty() // Return true if all opening tags are closed
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_valid_xml() {
        assert!(is_valid_xml("<Design><Code>hello world</Code></Design>"));
        assert!(is_valid_xml("<People></People>"));
        assert!(is_valid_xml("<People><Design><Code></Code></Design></People>"));
    }

    #[test]
    fn test_invalid_xml() {
        assert!(!is_valid_xml("tutorial"));
        assert!(!is_valid_xml("<tutorial>"));
        assert!(!is_valid_xml("t>>><<"));
        assert!(!is_valid_xml("<<tu>>torial"));
        assert!(!is_valid_xml("<tutorial><topic>"));
        assert!(!is_valid_xml("<Design><Code>hello world</Code></Design><People>"));
        assert!(!is_valid_xml("<People><Design><Code>hello world</People></Code></Design>"));
        assert!(!is_valid_xml("<People age=”1”>hello world</People>"));
        assert!(!is_valid_xml("<People date='<>01/01/2000'>hello world</People>"));
    }
}
