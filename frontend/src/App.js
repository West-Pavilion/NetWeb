import React, { useState } from 'react';
import './App.css';

function App() {
  const [command, setCommand] = useState('curl');
  const [url, setUrl] = useState('');
  const [customCommand, setCustomCommand] = useState('');
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const response = await fetch('/api/test', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          command: command,
          url: url,
          custom: customCommand,
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      setResult(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const formatOutput = (output) => {
    if (!output) return '';
    return output;
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>🌐 NetWeb - Network Testing Tool</h1>
        <p>Test network connectivity using various commands</p>
      </header>

      <main className="App-main">
        <form onSubmit={handleSubmit} className="test-form">
          <div className="form-group">
            <label htmlFor="command">Command Type:</label>
            <select
              id="command"
              value={command}
              onChange={(e) => setCommand(e.target.value)}
              className="form-control"
            >
              <option value="curl">cURL - HTTP Request</option>
              <option value="ping">Ping - ICMP Echo</option>
              <option value="tracert">Traceroute - Path Tracing</option>
              <option value="custom">Custom Command</option>
            </select>
          </div>

          {command === 'custom' && (
            <div className="form-group">
              <label htmlFor="customCommand">Custom Command:</label>
              <input
                type="text"
                id="customCommand"
                value={customCommand}
                onChange={(e) => setCustomCommand(e.target.value)}
                placeholder="e.g., nslookup {url} or dig {url}"
                className="form-control"
              />
              <small className="form-hint">Use {'{url}'} as placeholder for the URL</small>
            </div>
          )}

          <div className="form-group">
            <label htmlFor="url">Target URL/Host:</label>
            <input
              type="text"
              id="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="e.g., https://example.com or example.com"
              className="form-control"
              required
            />
          </div>

          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Testing...' : 'Run Test'}
          </button>
        </form>

        {error && (
          <div className="error-box">
            <h3>Error</h3>
            <p>{error}</p>
          </div>
        )}

        {result && (
          <div className="results-container">
            <div className={`status-banner ${result.success ? 'success' : 'failure'}`}>
              {result.success ? '✓ Test Successful' : '✗ Test Failed'}
            </div>

            <div className="info-section">
              <h3>Connection Information</h3>
              <div className="info-grid">
                <div className="info-item">
                  <span className="info-label">Target:</span>
                  <span className="info-value">{result.connection.target}</span>
                </div>
                <div className="info-item">
                  <span className="info-label">Command:</span>
                  <span className="info-value">{result.command}</span>
                </div>
                <div className="info-item">
                  <span className="info-label">Duration:</span>
                  <span className="info-value">{result.duration}</span>
                </div>
                <div className="info-item">
                  <span className="info-label">OS:</span>
                  <span className="info-value">{result.connection.os}</span>
                </div>
                <div className="info-item">
                  <span className="info-label">Timestamp:</span>
                  <span className="info-value">{new Date(result.connection.timestamp).toLocaleString()}</span>
                </div>
              </div>
            </div>

            {result.error && (
              <div className="error-section">
                <h3>Error Details</h3>
                <pre className="output-box error-output">{result.error}</pre>
              </div>
            )}

            <div className="output-section">
              <h3>Response Output</h3>
              <pre className="output-box">{formatOutput(result.output)}</pre>
            </div>

            {result.metadata && Object.keys(result.metadata).length > 0 && (
              <div className="metadata-section">
                <h3>Metadata</h3>
                <div className="info-grid">
                  {Object.entries(result.metadata).map(([key, value]) => (
                    <div key={key} className="info-item">
                      <span className="info-label">{key}:</span>
                      <span className="info-value">{value}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
      </main>

      <footer className="App-footer">
        <p>Built with Go and React | Network Debugging Tool</p>
      </footer>
    </div>
  );
}

export default App;