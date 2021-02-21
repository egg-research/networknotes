import React, { useState, useEffect } from 'react';

import { AutoComplete, Input, Tag } from 'antd';
import './SearchBar.css';

export default function SearchBar({
  documents = [],
  keywords = [],
  newKeyword,
  selectDocument,
  selectKeyword,
  createKeyword,
  placeholder = 'Search',
}) {
  const initDocs = documents;
  const initKeywords = keywords;

  const [searchedDocs, setSearchedDocs] = useState(initDocs);
  const [searchedKeywords, setSearchedKeywords] = useState(initKeywords);
  const [value, setValue] = useState('');

  useEffect(() => {
    setSearchedDocs(documents);
    setSearchedKeywords(keywords);
    setValue('');
  }, [document, keywords]);

  const autocompleteOptions = [];
  if (keywords.length !== 0) {
    autocompleteOptions.push({
      label: <span>Keywords</span>,
      options: searchedKeywords.map((keyword) => ({
        value: keyword.id,
        label: <Tag key={keyword.id}>{keyword.name}</Tag>,
        type: 'keyword',
        data: keyword,
      })),
    });
  }

  if (documents.length !== 0) {
    autocompleteOptions.unshift({
      label: <span>Documents</span>,
      options: searchedDocs.map((document) => ({
        data: document,
        value: document.id,
        label: <div key={document.id}>{document.title}</div>,
        type: 'document',
      })),
    });
  }

  if (newKeyword) {
    autocompleteOptions.push({
      label: <span>New Keyword</span>,
      options: [
        {
          type: 'add',
          label: <div>Add keyword: {value}</div>,
          value,
        },
      ],
    });
  }

  const onSearch = (searchText) => {
    console.log(searchText);
    if (!searchText) {
      setSearchedDocs(initDocs);
      setSearchedKeywords(initKeywords);
      setValue('');
      return;
    }

    const newSearchedKeywords = keywords.filter((keyword) =>
      keyword.name?.includes(searchText)
    );

    console.log(documents);
    const newSearchedDocs = documents.filter((doc) =>
      doc.title?.includes(searchText)
    );

    setSearchedDocs(newSearchedDocs);
    setSearchedKeywords(newSearchedKeywords);
    setValue(searchText);
  };

  const onSelect = (_, instance) => {
    if (instance.type === 'document' && selectDocument) {
      selectDocument(instance.data);
    } else if (instance.type === 'keyword' && selectKeyword) {
      selectKeyword(instance.data);
    } else {
      createKeyword(value);
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
