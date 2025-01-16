import React, { useState } from 'react';
import getPathOnVersion from '../services/versions.jsx';

function ProcessForm({ ws, serverInfo }) {
    const [path, setPath] = useState(''); // Путь к приложению
    const [pid, setPid] = useState(''); // PID процесса для остановки
    const [edition, setEdition] = useState(''); // Редакция
    const [url, setUrl] = useState(''); // URL сервера

    const handleStart = () => {
        if (ws) {
            const pathTo = getPathOnVersion(serverInfo.edition, serverInfo.version)
            if (!pathTo) {
                return;
            }
            ws.send(JSON.stringify({ action: 'start', path: pathTo }));
        }
    };

    const handleStop = () => {
        if (ws && pid) {
            ws.send(JSON.stringify({ action: 'stop', pid }));
            setPid('');
        }
    };

    const handleUpdateXMLConfig = () => {
        if (ws && edition && url) {
            ws.send(JSON.stringify({
                action: 'update_xml_config',
                edition,
                url,
            }));
            setEdition('');
            setUrl('');
        }
    };

    return (
        <div className="process-form">
            <h2>Manage Processes</h2>
            <div>
                <label>
                    Path to executable:
                    <input
                        type="text"
                        value={path}
                        onChange={(e) => setPath(e.target.value)}
                        placeholder="Enter path to executable"
                    />
                </label>
                <button onClick={handleStart}>Start</button>
            </div>
            <div>
                <label>
                    PID to stop:
                    <input
                        type="text"
                        value={pid}
                        onChange={(e) => setPid(e.target.value)}
                        placeholder="Enter PID"
                    />
                </label>
                <button onClick={handleStop}>Stop</button>
            </div>
            <div>
                <h3>Update XML Config</h3>
                <label>
                    Edition:
                    <input
                        type="text"
                        value={edition}
                        onChange={(e) => setEdition(e.target.value)}
                        placeholder="Enter edition"
                    />
                </label>
                <label>
                    Server URL:
                    <input
                        type="text"
                        value={url}
                        onChange={(e) => setUrl(e.target.value)}
                        placeholder="Enter server URL"
                    />
                </label>
                <button onClick={handleUpdateXMLConfig}>Update XML Config</button>
            </div>
        </div>
    );
}

export default ProcessForm;