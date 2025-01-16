import React, { useState, useEffect } from 'react';
import ProcessList from './components/ProcessList';
import ProcessForm from './components/ProcessForm';
import './App.css';
import ServerInfo from "./components/ServerInfo.jsx";
import ServersList from "./components/ServersList.jsx";
import {inConstructor} from "eslint-plugin-react/lib/util/ast.js";


//C:\Program Files (x86)\AnyDesk\AnyDesk.exe
function App() {
    const [processes, setProcesses] = useState([]); // Состояние для хранения процессов
    const [servers, setServers] = useState([]);
    const [serverInfo, setServerInfo] = useState(null);
    const [url, setUrl] = useState('');
    const [input, setInput] = useState('');
    const [infoLoading, setInfoLoading] = useState(false);
    const [buttonEnable, setButtonEnable] = useState(false);

    const [ws, setWs] = useState(null); // Состояние для WebSocket


    // Подключение к WebSocket
    useEffect(() => {
        const websocket = new WebSocket('ws://localhost:4000/ws');

        websocket.addEventListener('open', () => {
            console.log('WebSocket соединение установлено');
            websocket.send(JSON.stringify({ action: 'get_processes' }));
        });

        websocket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            console.log(data);
            if (data.action === 'server_info') {
                // Обновляем состояние с информацией о сервере
                setServerInfo({
                    edition: data.edition,
                    version: data.version,
                    state: data.state,
                });
                setInfoLoading(false)
            }
            if (data.action === 'start') {
                // Добавляем новый процесс
                setProcesses((prevProcesses) => [
                    ...prevProcesses,
                    { path: data.path, pid: data.pid, status: data.status, version: data.version, edition: data.edition, urlParam: data.urlParam },
                ]);
            } else if (data.action === 'stop') {
                // Обновляем статус процесса по pid
                setProcesses((prevProcesses) =>
                    prevProcesses.map((process) =>
                        process.pid === data.pid
                            ? { ...process, status: data.status } // Обновляем статус
                            : process // Оставляем без изменений
                    )
                );
            }
        };
        setWs(websocket);

        return () => websocket.close(); // Закрытие соединения при размонтировании
    }, []);

    return (
        <div className="App">

            <div className="split">
                <ServersList servers={servers}
                             setServers={setServers}
                             setUrl={setUrl} />
                <div className="server-panel">
                    <ServerInfo ws={ws}
                                serverInfo={serverInfo}
                                infoLoading={infoLoading}
                                setInfoLoading={setInfoLoading}
                                url={url}
                                setUrl={setUrl}
                                buttonEnable={buttonEnable}
                                setButtonEnable={setButtonEnable}
                                servers={servers}
                                setServers={setServers}
                                input={input}
                                setInput={setInput}
                                processes={processes}/>
                </div>

            <ProcessList processes={processes}
                         ws={ws}
                         serverInfo={serverInfo}
                         url={url}
                         buttonEnable={buttonEnable}
                         setInput={setInput}
                         setProcesses={setProcesses}/>
            </div>

        </div>
    );
}

export default App;