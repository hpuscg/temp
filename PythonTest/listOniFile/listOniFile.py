# _*_ coding:utf-8 _*_
import os

def main():
    path = 'D:\\oni\\' # oni的tar包路径
    new_path = "D:\\oni\\new\\" # 解压后的oni的tar包存储路径
    class_path = "D:\\oni\\class\\" # 分类后的oni的tar包存储路径
    # 处理tar包并进行分类
    handler(path, new_path, class_path)

def handler(path, new_path, class_path):
    # 获取tar包的文件列表
    tar_oni_file_list = get_file_dir_list(path)
    for tar_file in tar_oni_file_list:
        if not tar_file.endswith("tar"):
            continue
        print(tar_file)
        # oni_dir = path + tar_file
        # 解压tar包
        # dir_names = tar_file.split(".")
        # new_oni_dir = path +dir_names[0]
        # if not os.path.exists(new_oni_dir):
            # os.system("mkdir " + new_oni_dir)
        tar_oni_file(tar_file)
        os.system("mv *.tar.gz " + new_path)
    oni_file_list = get_file_dir_list(new_path)
    for oni_file_name in oni_file_list:
        print(oni_file_name)
        dir_names = oni_file_name.split("_")
        class_dir = class_path + dir_names[0]
        # 创建以uid为名的文件夹
        if not os.path.exists(class_dir):
            os.system("mkdir " + class_dir)
        os.system("mv " + new_path + oni_file_name + " " + class_dir)
    '''
    # 获取所有的tar包的文件名
    new_oni_file_list = get_file_dir_list(new_path)
    for new_oni_file in new_oni_file_list:
        print(new_oni_file)
        new_oni_file_path = new_path + new_oni_file
        dir_names = new_oni_file.split("_")
        class_dir = class_path + dir_names[0]
        # 创建以uid为名的文件夹
        if not os.path.exists(class_dir):
            os.system("mkdir " + class_dir)
        # 归类oni到每个uid下
        os.system("mv " + new_oni_file_path + " " + class_dir)
    '''


# 获取文件下的目录
def get_file_dir_list(path):
    file_dir_list = os.listdir(path)
    return file_dir_list

# 解压oni的tar包文件
def tar_oni_file(filename):
    str = 'tar -zxvf ' + filename
    print(str)
    os.system(str)

# 测试函数
def test_parm(path):
    print(path)
    return path, 2, 4

if __name__ == '__main__':
    main()
