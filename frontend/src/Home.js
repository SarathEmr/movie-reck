import React, { useState } from 'react';

function Home({ user, onLogout }) {
  const [selectedGenres, setSelectedGenres] = useState([]);
  const [recommendations, setRecommendations] = useState([]);

  const genres = ['sci-fi', 'drama', 'romcom', 'horror', 'horror-comedy', 'thriller'];

  const handleGenreChange = (genre) => {
    setSelectedGenres(prev =>
      prev.includes(genre)
        ? prev.filter(g => g !== genre)
        : [...prev, genre]
    );
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const queryParams = new URLSearchParams();
    queryParams.append('language', 'english');
    selectedGenres.forEach(genre => queryParams.append('genres', genre));

    fetch(`http://localhost:8080/movie/recommendation?${queryParams.toString()}`)
      .then(response => response.json())
      .then(data => setRecommendations(data.movies || []))
      .catch(err => console.error('Error fetching recommendations:', err));
  };

  return (
    <div>
      <h2>Welcome, {user.username}!</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Select your preferred genres:</label>
          {genres.map(genre => (
            <label key={genre} style={{ display: 'block', margin: '5px 0' }}>
              <input
                type="checkbox"
                checked={selectedGenres.includes(genre)}
                onChange={() => handleGenreChange(genre)}
              />
              {genre}
            </label>
          ))}
        </div>
        <button type="submit">Get Recommendations</button>
      </form>
      <h3>Movie Recommendations</h3>
      <ul>
        {recommendations.map((movie, index) => (
          <li key={movie.id || index}>
            {movie.name} ({movie.year}) - Genres: {movie.genre.join(', ')} - Rating: {movie.imdb_rating}
          </li>
        ))}
      </ul>
      <button onClick={onLogout}>Logout</button>
    </div>
  );
}

export default Home;