import React, { useState, useEffect } from 'react';

import { AutoComplete, Input, Tag } from 'antd';
import './SearchBar.css';

export default function SearchBar({
  documents = [],
  keywords = [],
  selectDocument,
  selectKeyword,
  placeholder = 'Search',
}) {
  // const defaultLen = 3;
  // const initDocs = documents.slice(
  //   0,
  //   documents.length > defaultLen ? defaultLen : documents.length
  // );

  // const initKeywords = keywords.slice(
  //   0,
  //   keywords.length > defaultLen ? defaultLen : keywords.length
  // );

  const initDocs = documents;
  const initKeywords = keywords;

  const [searchedDocs, setSearchedDocs] = useState(initDocs);
  const [searchedKeywords, setSearchedKeywords] = useState(initKeywords);
  const [value, setValue] = useState('');

  useEffect(() => {
    // const d = documents.slice(
    //   0,
    //   documents.length > defaultLen ? defaultLen : documents.length
    // );

    // const k = keywords.slice(
    //   0,
    //   keywords.length > defaultLen ? defaultLen : keywords.length
    // );
    setSearchedDocs(documents);
    setSearchedKeywords(keywords);
    setValue('');
  }, [document, keywords]);

  const autocompleteOptions = [
    {
      label: <span>Documents</span>,
      options: searchedDocs.map((document) => ({
        value: document,
        label: <div>{document}</div>,
        type: 'document',
      })),
    },
    {
      label: <span>Keywords</span>,
      options: searchedKeywords.map((keyword) => ({
        value: keyword,
        label: <Tag>{keyword}</Tag>,
        type: 'keyword',
      })),
    },
  ];

  const onSearch = (searchText) => {
    if (!searchText) {
      setSearchedDocs(initDocs);
      setSearchedKeywords(initKeywords);
      setValue('');
      return;
    }

    const newSearchedKeywords = keywords.filter((keyword) =>
      keyword.includes(searchText)
    );

    const newSearchedDocs = documents.filter((doc) => doc.includes(searchText));

    setSearchedDocs(newSearchedDocs);
    setSearchedKeywords(newSearchedKeywords);
    setValue(searchText);
  };

  const onSelect = (_, instance) => {
    if (instance.type === 'document' && selectDocument) {
      selectDocument(instance.value);
    } else if (instance.type === 'keyword' && selectKeyword) {
      selectKeyword(instance.value);
    }
    setValue('');
  };

  return (
    <AutoComplete
      options={autocompleteOptions}
      dropdownClassName='dropdown-category'
      style={{ width: '100%' }}
      value={value}
      onSearch={onSearch}
      onSelect={onSelect}
    >
      <Input.Search placeholder={placeholder} />
    </AutoComplete>
  );
}
