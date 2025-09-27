import { useState, useEffect } from "react";
import Autosuggest from "react-autosuggest";
import "./App.css";
import names from "./assets/tech_terms_flat_updated.json";

const terms = names;

function App() {
  const [value, setValue] = useState("");
  const [suggestions, setSuggestions] = useState([]);
  const [interests, setInterests] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  const [selectedCategory, setSelectedCategory] = useState("new");
  const [filteredNews, setFilteredNews] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const newsPerPage = 20;

  // Load interests from localStorage
  useEffect(() => {
    try {
      const stored = localStorage.getItem("user_interests");
      if (stored) {
        const data = JSON.parse(stored);

        const names = data.Names || data.names || [];
        if (Array.isArray(names)) {
          const loadedInterests = names.map((name) => ({ name }));
          setInterests(loadedInterests);
        }
      }
    } catch (err) {
      console.error("Failed to parse localStorage:", err);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchNews = async () => {
    try {
      console.log("Fetching news for category:", selectedCategory);

      const saved = localStorage.getItem("user_interests");
      if (!saved) {
        console.log("No saved interests, showing empty news");
        setFilteredNews([]);
        return;
      }

      const parsed = JSON.parse(saved);
      console.log("Sending to backend:", parsed);

      const res = await fetch("http://localhost:8080/getNews", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(parsed),
      });

      if (!res.ok) {
        console.error("Network response not ok:", res.status);
        setFilteredNews([]);
        return;
      }

      const data = await res.json();
      console.log("Got news data:", data);

      setFilteredNews(
        selectedCategory === "new"
          ? data.newStories || []
          : data.topStories || [],
      );

      setCurrentPage(1);
    } catch (err) {
      console.error("Fetch error:", err);
      // Don't show error to user, just show empty news
      setFilteredNews([]);
    }
  };

  useEffect(() => {
    if (isLoading) return;

    fetchNews();
  }, [selectedCategory, isLoading]);

  const getSuggestions = (value) => {
    const inputValue = value.trim().toLowerCase();
    const inputLength = inputValue.length;

    return inputLength === 0
      ? []
      : terms.filter(
          (term) =>
            term.name.toLowerCase().slice(0, inputLength) === inputValue &&
            !interests.some((interest) => interest.name === term.name),
        );
  };

  const getSuggestionValue = (suggestion) => suggestion.name;
  const renderSuggestion = (suggestion) => <div>{suggestion.name}</div>;

  const onChange = (event, { newValue }) => {
    setValue(newValue);
  };

  const onSuggestionsFetchRequested = ({ value }) => {
    setSuggestions(getSuggestions(value));
  };

  const onSuggestionsClearRequested = () => {
    setSuggestions([]);
  };

  const onSuggestionSelected = (event, { suggestion }) => {
    if (!interests.some((interest) => interest.name === suggestion.name)) {
      setInterests([...interests, suggestion]);
    }
    setValue("");
  };

  const removeInterest = (name) => {
    setInterests(interests.filter((interest) => interest.name !== name));
  };

  const handleSubmit = async () => {
    const payload = interests.map((interest) => interest.name);

    try {
      console.log("Submitting interests:", payload);

      const res = await fetch("http://localhost:8080/saveInterests", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ names: payload }),
      });

      if (!res.ok) {
        console.error("Submit response not ok:", res.status);
        return;
      }

      const data = await res.json();
      console.log("Submit response:", data);

      // Make sure it's in a consistent shape
      const formattedData = {
        Names: data.Names || [],
        Embeddings: data.Embeddings || [],
      };

      // Save full response (Names + Embeddings) to localStorage
      localStorage.setItem("user_interests", JSON.stringify(formattedData));
      console.log("Saved to localStorage:", formattedData);

      if (formattedData.Names.length > 0) {
        const formatted = formattedData.Names.map((name) => ({ name }));
        setInterests(formatted);
      }
      await fetchNews(data);
    } catch (err) {
      console.error("Error submitting interests:", err);
    }
  };

  const handleCategoryChange = (category) => {
    setSelectedCategory(category);
  };

  const indexOfLastNews = currentPage * newsPerPage;
  const indexOfFirstNews = indexOfLastNews - newsPerPage;
  const currentNews = filteredNews.slice(indexOfFirstNews, indexOfLastNews);
  const totalPages = Math.ceil(filteredNews.length / newsPerPage);

  const inputProps = {
    placeholder: "Add Interests",
    value,
    onChange,
    style: { fontSize: "18px", padding: "12px 16px", height: "50px" },
  };

  console.log("Render - isLoading:", isLoading, "interests:", interests);

  if (isLoading) {
    return (
      <div className="w-full min-h-screen bg-[#faedcd] custom-font flex flex-col justify-center items-center">
        <div className="text-[#d4a373] custom-font text-5xl mb-10">GO-HN</div>
        <div className="text-[#d4a373] text-xl">Loading your interests...</div>
      </div>
    );
  }

  return (
    <div className="w-full min-h-screen bg-[#faedcd] custom-font flex flex-col justify-start items-center p-4">
      <div className="text-[#d4a373] custom-font text-5xl mb-10">GO-HN</div>

      <Autosuggest
        suggestions={suggestions}
        onSuggestionsFetchRequested={onSuggestionsFetchRequested}
        onSuggestionsClearRequested={onSuggestionsClearRequested}
        getSuggestionValue={getSuggestionValue}
        renderSuggestion={renderSuggestion}
        onSuggestionSelected={onSuggestionSelected}
        inputProps={inputProps}
      />

      <div className="mt-4 flex flex-wrap gap-2">
        {interests
          .filter((interest) => interest && interest.name) // 👈 Add this filter for safety
          .map((interest) => (
            <div
              key={interest.name}
              className="bg-[#d4a373] text-white px-3 py-1 rounded-full flex items-center"
            >
              <span>{interest.name}</span>
              <button
                onClick={() => removeInterest(interest.name)}
                className="ml-2 text-sm hover:text-gray-300"
              >
                ×
              </button>
            </div>
          ))}
      </div>

      <button
        onClick={handleSubmit}
        className="mt-4 bg-[#d4a373] text-white px-4 py-2 rounded hover:bg-[#c79a5b]"
      >
        Submit Interests
      </button>

      <div className="w-full max-w-md mt-4">
        <select
          value={selectedCategory}
          onChange={(e) => handleCategoryChange(e.target.value)}
          className="w-full p-3 border border-[#d4a373] rounded bg-white text-[#d4a373] font-medium"
        >
          <option value="new">New</option>
          <option value="top">Top</option>
        </select>
      </div>

      <div className="w-full max-w-4xl mt-8">
        <h2 className="text-2xl text-[#d4a373] mb-6 text-center">
          {selectedCategory === "new" ? "Latest News" : "Top News"}
        </h2>

        <div className="grid gap-4">
          {currentNews.length === 0 ? (
            <div className="text-center text-[#d4a373] py-8">
              {filteredNews.length === 0
                ? "No news available. Try submitting some interests first!"
                : "No news found."}
            </div>
          ) : (
            currentNews
              .filter((news) => news && news.id && news.title && news.url)
              .map((news) => (
                <div
                  key={news.id}
                  className="bg-white p-6 rounded-lg shadow-md border border-[#d4a373]/20"
                >
                  <h3 className="text-xl font-semibold text-[#d4a373] mb-2">
                    {news.title}
                  </h3>
                  <p className="text-[#d4a373] mb-3">{news.text}</p>
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-[#d4a373]">by {news.by}</span>
                    <a
                      href={news.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="bg-[#d4a373] text-white px-4 py-2 rounded text-sm hover:bg-[#c79a5b]"
                    >
                      Read Full Article
                    </a>
                  </div>
                </div>
              ))
          )}
        </div>

        {totalPages > 1 && (
          <div className="flex justify-center items-center mt-8 gap-2">
            <button
              onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
              disabled={currentPage === 1}
              className="px-4 py-2 bg-[#d4a373] text-white rounded disabled:bg-gray-300 disabled:cursor-not-allowed"
            >
              Previous
            </button>

            <span className="px-4 py-2 text-[#d4a373]">
              Page {currentPage} of {totalPages}
            </span>

            <button
              onClick={() =>
                setCurrentPage((prev) => Math.min(prev + 1, totalPages))
              }
              disabled={currentPage === totalPages}
              className="px-4 py-2 bg-[#d4a373] text-white rounded disabled:bg-gray-300 disabled:cursor-not-allowed"
            >
              Next
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
