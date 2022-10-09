import yaml

# yaml file content
# a:
# - a1
# - a2
# b:
#   b1:
#     b11: true
#     b12: helm
#   b2:
#     b21: false
#     b22: kustomize
file_path = "./file.yaml"

def read_file(file_path):
    with open(file_path, 'r') as f:
        file_content = yaml.safe_load(f)
    return file_content

def update_and_write_file(old_file):
    file_content = read_file(old_file)
    file_content['new_key'] = 'value' # add key in yaml, ok even key doesn't exist 
    
    # write to file
    with open(file_path, 'w', encoding='utf-8') as f:
        yaml.dump(file_content, f)


file_content = read_file(file_path)   

print(file_content)
# {'a': ['a1', 'a2'], 'b': {'b1': {'b11': True, 'b12': 'helm'}, 'b2': {'b21': False, 'b22': 'kustomize'}}}
# {'b1': {'b11': True, 'b12': 'helm'}, 'b2': {'b21': False, 'b22': 'kustomize'}}

print(file_content['b']) 
# {'b1': {'b11': True, 'b12': 'helm'}, 'b2': {'b21': False, 'b22': 'kustomize'}}

# check if range works for multiple yaml objects, not array
for ele in file_content['b']:
    print(file_content['b'][ele])
# {'b11': True, 'b12': 'helm'}
# {'b21': False, 'b22': 'kustomize'}

update_and_write_file(file_path)
