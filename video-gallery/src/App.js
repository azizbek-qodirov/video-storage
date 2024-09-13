import React, { useCallback, useEffect, useState } from 'react';

const VideoPlayer = ({ video }) => {
  const handleDoubleTap = useCallback((e) => {
    const videoElement = e.target;
    const now = new Date().getTime();
    const timeSince = now - (videoElement.lastTap || 0);
    
    if (timeSince < 300 && timeSince > 0) {
      if (e.nativeEvent.clientX < videoElement.offsetWidth / 2) {
        videoElement.currentTime = Math.max(videoElement.currentTime - 5, 0);
      } else {
        videoElement.currentTime = Math.min(videoElement.currentTime + 5, videoElement.duration);
      }
    }
    
    videoElement.lastTap = now;
  }, []);

  return (
    <div className="video-container">
      <video 
        src={`http://${video.url}`}
        controls
        onClick={handleDoubleTap}
      />
      <h3>{video.name}</h3>
    </div>
  );
};

const VideoGallery = () => {
  const [videos, setVideos] = useState([]);
  const [error, setError] = useState(null);

  const isVideoFile = useCallback((filename) => {
    const videoExtensions = ['.mp4', '.webm', '.ogg', '.mov'];
    return videoExtensions.some(ext => filename.toLowerCase().endsWith(ext));
  }, []);

  useEffect(() => {
    const fetchVideos = async () => {
      try {
        const response = await fetch('http://localhost:8088/api/v1/videos');
        const files = await response.json();
        const filteredVideos = files.filter(file => isVideoFile(file.name));
        setVideos(filteredVideos);
      } catch (error) {
        console.error('Error fetching videos:', error);
        setError('Error loading videos: ' + error.message);
      }
    };

    fetchVideos();
  }, [isVideoFile]);

  if (error) {
    return <p>{error}</p>;
  }

  if (videos.length === 0) {
    return <p>No videos found.</p>;
  }

  return (
    <div className="video-grid">
      {videos.map(video => (
        <VideoPlayer key={video.id} video={video} />
      ))}
    </div>
  );
};

const App = () => {
  return (
    <div className="app">
      <h1>Enhanced Video Gallery</h1>
      <VideoGallery />
    </div>
  );
};

// Styles
const styles = `
  body { 
    font-family: Arial, sans-serif; 
    margin: 0;
    padding: 20px;
    background-color: #f0f0f0;
  }
  h1 {
    color: #333;
    text-align: center;
  }
  .video-grid { 
    display: flex; 
    flex-wrap: wrap; 
    gap: 20px; 
    justify-content: center;
  }
  .video-container { 
    width: 320px; 
    margin-bottom: 20px; 
    background-color: white;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 4px 6px rgba(0,0,0,0.1);
  }
  .video-container video {
    width: 100%;
    display: block;
  }
  .video-container h3 {
    margin: 10px;
    font-size: 16px;
    color: #333;
  }
`;

// Add styles to the document
const styleElement = document.createElement('style');
styleElement.textContent = styles;
document.head.appendChild(styleElement);

export default App;
