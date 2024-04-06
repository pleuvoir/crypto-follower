const { app, BrowserWindow, Menu, globalShortcut } = require("electron");
const openAboutWindow = require("about-window").default;
const freePort = require("freeport");
const fetch = require("electron-fetch").default;

const { execFile } = require("child_process");

let backendPort = 8000;
let win = null;

//添加了一个命令行开关，用于禁用 Electron 窗口中的 Web 安全性功能，通常用于在开发过程中解决跨域访问的问题。
app.commandLine.appendSwitch("disable-web-security");
app.whenReady().then(() => {
    // 在应用程序准备就绪时执行一些操作

    freePort(function (err, port) {
        // 通过 free port 模块获取一个空闲端口

       // 定义了一个变量 backendBin，用于存储后端服务的可执行文件的路径。
        let backendBin = "./crypto-follower_backend";

        let binFolder = __dirname + "/../bin";
        console.info(binFolder)

        // 使用 execFile 方法执行后端服务的可执行文件，传入了命令行参数 --port 和获取到的空闲端口号，并设置了执行文件的工作目录和环境变量。
        try {
            const child = execFile(backendBin, ["--port", port], {
                cwd: binFolder
            });
            let url = "http://127.0.0.1:" + backendPort + "/app";

            child.stdout.on("data", (data) => {

                // 当后端服务输出数据时执行的回调函数
                console.info(data)


                //检查窗口对象是否已经存在，如果存在则直接返回，避免重复创建窗口。
                if (win != null) {
                    return;
                }

                    win = new BrowserWindow({
                    title: "crypto-follower",
                    maximizable: true,
                    show: false,
                    resizable: true,
                    webPreferences: {
                        nodeIntegration: true,
                        enableRemoteModule: true,
                        webSecurity: false,
                    },
                });

                setTimeout(() => {
                    win.loadURL(url, {
                        userAgent: "App",
                    });

                    win.webContents.openDevTools();

                    globalShortcut.register("f5", function () {
                        win.reload();
                    });
                    globalShortcut.register("CommandOrControl+R", function () {
                        win.reload();
                    });
                    win.maximize();
                    win.show();
                }, 1000);
            }
            );

            const appMenu = {
                label: "app",
                role: "appMenu",
            };

            const editMenu = {
                label: "编辑",
                submenu: [
                    {
                        label: "Undo",
                        accelerator: "CmdOrCtrl+Z",
                        role: "undo",
                    },
                    {
                        label: "Redo",
                        accelerator: "Shift+CmdOrCtrl+Z",
                        role: "redo",
                    },
                    {
                        type: "separator",
                    },
                    {
                        label: "Cut",
                        accelerator: "CmdOrCtrl+X",
                        role: "cut",
                    },
                    {
                        label: "Copy",
                        accelerator: "CmdOrCtrl+C",
                        role: "copy",
                    },
                    {
                        label: "Paste",
                        accelerator: "CmdOrCtrl+V",
                        role: "paste",
                    },
                    {
                        label: "Select All",
                        accelerator: "CmdOrCtrl+A",
                        role: "selectAll",
                    },
                ],
            };

            const menu = Menu.buildFromTemplate([appMenu, editMenu]);

            Menu.setApplicationMenu(menu);
        } catch (e) {
            console.error(e)
        }
    });
});

app.on("before-quit", function () {
    if (process.env.NODE_ENV !== "dev") {
        new Promise((resolve, reject) => {
            fetch("http://127.0.0.1:" + backendPort + "/rpc/Quit", {
                method: "get",
            }).then(function (response) {
                resolve(response);
            }),
                (error) => {
                    reject(new Error(error.message));
                };
        });
    }
});