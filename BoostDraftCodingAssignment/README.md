A Simple XMl validator.
## Usage
```shell
# install
cargo install --path .

# directly use as a command tool, remember to add /.cargo/bin to your PATH
xml_validator.exe "xml_string"
```

## Description
This project is a simple xml validator tool, which validates if an xml string is structured properly.
- Input: an XML string be passed as the first command line argument.
- Output: `Valid` if given XML is valid, otherwise `Invalid` with a new line

Here's an overview of the changes:
- Main Program (main.rs):
The main program handles command line arguments to receive an XML string.
It checks for the presence of the XML string as a command line argument and prints an error message if absent or if too many arguments are provided.
The XML string is passed to the is_valid_xml function for validation.
Based on the validation result, the program outputs either "Valid" or "Invalid".
- XML Validation Function (is_valid_xml):
The is_valid_xml function checks if the given XML string is valid according to specified rules. 
  - return false if string did not contains any "<" or ">" 
  - use a stack to store xml tags
  - iterate over the string, when encountering opening tag (<xml_tag> without starting '/' ) push it to stack, when encountering a closed tag (eg: </xml_tag>), compare the tag string with the stack top, if they didn't match then return `false`, else pop out top element from stack
  - after iterating whole string, if stack is empty which means all opening tags are matched with closed tags, then return `true` else return `false`



Tests evidence:
```rust
running 2 tests
test tests::test_invalid_xml ... ok
test tests::test_valid_xml ... ok

test result: ok. 2 passed; 0 failed; 0 ignored; 0 measured; 0 filtered out; finished in 0.00s
```