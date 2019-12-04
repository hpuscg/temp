import os


def main():
    deleteDockerImages()


def deleteDockerImages():
    imageList = os.popen("sudo docker images").readlines()
    for image in imageList:
        print(image)
        tar = image.split(" ", 1)[1].strip().split(" ", 1)[1].strip().split(" ", 1)[0]
        print(tar)
        """name = image.split(" ", 1)[0]
        tag = image.split(" ", 1)[1].strip().split(" ", 1)[0]
        if "TAG" == tag :
            continue
        print(tag)
        name = image.split(" ", 1)[0].split("/")[1].split("-")[1]
        print(name)"""
        # if name != "REPOSITORY" :
        #    if name.split("/")[1] != "armhf-pioneer" and name.split("/")[1] != "armhf-libra-init":
        #        os.system("sudo docker rmi -f " + tar)
        # if "<none>" == name:
        #   print("===22222====", tar)
        os.system("sudo docker rmi -f " + tar)


if __name__ == '__main__':
    main()
