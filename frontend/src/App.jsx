import React, { useState } from "react";
import Autosuggest from "react-autosuggest";
import "./App.css";

// List of suggestions
const languages = [
  { name: "C", year: 1972 },
  { name: "Elm", year: 2012 },
  { name: "Go", year: 2009 },
  { name: "Java", year: 1995 },
  { name: "JavaScript", year: 1995 },
  { name: "Python", year: 1991 },
  { name: "Rust", year: 2010 },
];

// Suggestion logic
const getSuggestions = (value) => {
  const inputValue = value.trim().toLowerCase();
  const inputLength = inputValue.length;

  return inputLength === 0
    ? []
    : languages.filter(
        (lang) =>
          lang.name.toLowerCase().slice(0, inputLength) === inputValue &&
          !selectedInterests.some((interest) => interest.name === lang.name),
      );
};

const getSuggestionValue = (suggestion) => suggestion.name;

const renderSuggestion = (suggestion) => <div>{suggestion.name}</div>;

// This will be updated inside the component
let selectedInterests = [];

function App() {
  const [value, setValue] = useState("");
  const [suggestions, setSuggestions] = useState([]);
  const [interests, setInterests] = useState([]);

  // Keep reference up-to-date for filtering suggestions
  selectedInterests = interests;

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
      const newInterests = [...interests, suggestion];
      setInterests(newInterests);
    }
    setValue("");
  };

  const removeInterest = (name) => {
    setInterests(interests.filter((interest) => interest.name !== name));
  };

  const handleSubmit = () => {
    // Send interests to backend
    const payload = interests.map((interest) => interest.name);
    console.log("Interests to send:", payload);

    // Example fetch request
    fetch("/api/interests", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ interests: payload }),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log("Server response:", data);
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  };

  const inputProps = {
    placeholder: "Type a programming language",
    value,
    onChange,
  };

  return (
    <div className="w-full min-h-screen bg-[#faedcd] custom-font flex flex-col justify-start items-center p-4">
      <div className="text-[#d4a373] custom-font text-5xl mb-10">GO-HN</div>

      <div className="w-full max-w-md">
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
          {interests.map((interest) => (
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
      </div>
    </div>
  );
}

export default App;
