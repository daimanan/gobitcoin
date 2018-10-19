# gobitcoin使用GO语言实现一个比特币系统

- 版本升级

  - 2019-10-19

    - 添加创建区块命令行

    ```bash
    createChain --address ADDRESS "创建一个区块"
    addBlock --data DATA "添加一个区块到区块链中"
    printChain           "打印所有区块信息"
    ```

  - 2018-10-18 

       - 添加windows测试用批处理

            run.bat   测试本地数据库，本地可执行程序，重新运行新系统

            bp.bat    快速打印区块链信息

       - 增加命令行模式

       ```bash
       addBlock --data DATA "添加一个区块到区块链中"
       printChain           "打印所有区块信息"
       ```

  - 2018-10-17 使用bolt本地数据改写区块链存放，并引入pow算力挖矿

  - 2018-10-16 实现最基础的比特币系统，用切片作为区块链在内存中存放区块

