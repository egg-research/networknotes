import React from 'react';

import { AutoComplete, Input, Tag } from 'antd';
import './SearchBar.css';

export default function SearchBar({ documents, keywords }) {
  const autocompleteOptions = [
    {
      label: <span>Documents</span>,
      options: documents.map((document) => ({
        value: document,
        label: <div>{document}</div>,
      })),
    },
    {
      label: <span>Keywords</span>,
      options: keywords.map((keyword) => ({
        value: keyword,
        label: <Tag>{keyword}</Tag>,
      })),
    },
  ];

  return (
    <AutoComplete
      options={autocompleteOptions}
      dropdownClassName='dropdown-category'
      style={{ width: 250 }}
    >
      <Input.Search placeholder='Search' />
    </AutoComplete>
  );
}
