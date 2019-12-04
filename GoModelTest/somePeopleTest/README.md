
## 测试程序使用说明

* 在Windows中双击即可启动测试程序

* 程序启动后，会保持一个黑窗口的状态，用于输出测试结果


## 配置文件说明：

* peoplenum是限定的人数下限，

* timelimitin是当室内小于限定人数时人员进入的时间间隔，

* timelimitout是当室内人数小于限定人数时人员出来的总用时，

* libraaddr是一体机的IP，

* ioaddr是继电器的IP，

* linein是进入方向的线ID，

* lineout是出来方向的线ID；


## 测试效果说明

### 测试效果展示方法

* 第一种展示方法是在程序启动的黑窗口种直接输出测试信息，信息全面，分类细致

* 第二种展示方法是通过开关量来展示

### 测试效果判断

* 通过开关量观察测试效果：开关量的指示灯亮起，间隔一秒有关掉，即为有事件触发

   ```
   1、Do1亮起，说明在timelimitin秒内，peoplenum人未同时进入

   2、Do2亮起，说明在timelimitout秒内，peoplenum人未同时出来
   ```
* 通过黑窗口观察测试效果：直接输出测试信息
    ```
    1、在timelimitin秒内，peoplenum人未同时进入信息：
    the people num is: ？, come in use time upper ？ s


    2、在timelimitout秒内，peoplenum人未同时出来信息：
    the home people num is: ？, lower ？


    3、实时进入和出来以及屋内总人数信息：
    total num is : ？,in num is : ？,out num is : ？
    ```
