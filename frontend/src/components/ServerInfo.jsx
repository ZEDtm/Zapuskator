import React, { useCallback, useEffect, useState } from 'react';
import getPathOnVersion from "../services/versions.jsx";
let counter = 0;
const ServerInfo = ({ ws, serverInfo, infoLoading, setInfoLoading, url, setUrl, buttonEnable, setButtonEnable, servers, setServers,
                        input, setInput, processes}) => {
    const [error, setError] = useState(null); // Сообщение об ошибке
    const [host, setHost] = useState(null);
    const [port, setPort] = useState(null);
    const [protocol, setProtocol] = useState(null);
    const [loading, setLoading] = useState(false);

    const cleanUrl = (url) => {
        let cleanedUrl = url.includes('/resto') ? url.split('/resto')[0] : url;
        cleanedUrl = cleanedUrl.endsWith('/') ? cleanedUrl.slice(0, -1) : cleanedUrl;
        return cleanedUrl;
    };

    const urlRegex = /(?:(?:https?|ftp):\/\/)?[\w/\-?=%.]+\.[\w/\-?=%.]+/gi;

    const handleGetServerInfo = useCallback(() => {
        if (ws && url) {
            setError(null);
            ws.send(JSON.stringify({ action: 'get_server_info', url: cleanUrl(url) }));
        }
    }, [ws, url]);

    const handleUrlChange = useCallback(() => {
        setInfoLoading(true);
        handleGetServerInfo();

    }, [handleGetServerInfo]);

    useEffect(() => {
        const debounceTimer = setTimeout(() => {
            if (url) {
                handleUrlChange(url);
            }
        }, 500); // Задержка 500 мс

        return () => clearTimeout(debounceTimer);
    }, [url, handleUrlChange]);

    const handleInputChange = useCallback((input) => {
        const protocolMatch = input.match(/^(https?:\/\/)/); // Префикс (http:// или https://)
        const domainMatch = input.match(/^https?:\/\/([^\/:]+)/); // Домен
        const portMatch = input.match(/:(\d+)/); // Порт

        const newProtocol = protocolMatch ? protocolMatch[1] : protocol ? protocol : portMatch? 'http://': 'https://';
        const newHost = domainMatch ? domainMatch[1] : input;
        const newPort = portMatch ? portMatch[1] : protocol === 'https://' ? null : port;
        if (port){
            const host = newHost.split(':'+port)[0]
            setProtocol(newProtocol);
            setHost(host);
            setPort(newPort);
            setInput(host);
        }else {
            setProtocol(newProtocol);
            setHost(newHost);
            setPort(newPort);
            setInput(newHost);
        }

    }, [protocol, port]);

    useEffect(() => {
        if (input) {
            handleInputChange(input);
        }
    }, [input, handleInputChange]);

    const changeProtocol = () => {
        if (protocol === 'https://') {
            setProtocol('http://');
        } else {
            setProtocol('https://');
            setPort(null);
        }
    };

    useEffect(() => {
        const debounceTimer = setTimeout(() => {
            if (protocol && host) {
                const newUrl = port ? `${protocol}${host}:${port}` : `${protocol}${host}`;
                setUrl(newUrl);
                setInput(host); // Обновляем input только если host изменился
            }
        }, 1000);

        return () => clearTimeout(debounceTimer);
    }, [protocol, host, port]);

    const handleButtonClick = () => {
        const newUrl = port ? `${protocol}${host}:${port}` : `${protocol}${host}`;
        const cleanedUrl = cleanUrl(newUrl);

        if (urlRegex.test(cleanedUrl)) {
            setInfoLoading(true);
            setUrl(cleanedUrl);
            handleGetServerInfo();
            const currentServers = servers.filter(server => server !== url);
            setServers([...currentServers, url]);
        } else {
            setError('Invalid URL format');
            setButtonEnable(false);
        }
    };

    const handleStart = () => {
        if (ws) {
            const path = getPathOnVersion(serverInfo.edition, serverInfo.version);
            if (!path) {
                return;
            }
            if (ws && serverInfo.edition && url) {
                ws.send(JSON.stringify({action: 'start', path: path, edition: serverInfo.edition, urlParam: url, version: serverInfo.version}));
            }
        }
    };

    useEffect(() => {
        if (cleanUrl(url)) {
            setButtonEnable(true);
        } else {
            setButtonEnable(false);
        }
    }, [url]);

    return (
        <>
            <h1><span className="word">Z</span>апускатор <span className="version">2.0</span></h1>
            <div className="server-info">
                <div>
                    <label>
                        {protocol && <button className="http-s-btn" onClick={changeProtocol}>{protocol}</button>}
                        <input
                            type="text"
                            value={input}
                            onChange={(e) => setInput(e.target.value)}
                            placeholder="URL"
                        />
                        {port && <button className="http-s-btn">{port}</button>}
                    </label>
                    <button disabled={!buttonEnable} onClick={handleButtonClick}>Получить</button>
                    <button disabled={!buttonEnable} onClick={handleStart}>Запустить</button>
                </div>

                {serverInfo && serverInfo.version && !infoLoading ? (
                    <div>
                        {error && <p style={{ color: 'red' }}>{error}</p>}
                        <>
                            <h2>Информация о сервере</h2>
                            <p>{url}</p>
                            <p>Редакция: {serverInfo.edition === 'default' ? 'Office' : 'Chain'}</p>
                            <p>Версия: {serverInfo.version}</p>
                            <p>Статус: {serverInfo.state}</p>
                        </>
                    </div>
                ) : infoLoading ? (
                    <div className="loader">
                        <div className="terminal-loader"></div>
                    </div>
                ) : (
                    <p>Просто вставь URL</p>
                )}
            </div>
        </>
    );
};

export default ServerInfo;