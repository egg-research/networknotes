import React from 'react';

import { Typography } from 'antd';

const { Text } = Typography;

export function HeavyText({ children }) {
  return <Text style={{ fontWeight: 'bold', fontSize: 14 }}>{children}</Text>;
}

export function WideHeavyText({ children }) {
  return (
    <Text style={{ fontWeight: 'bold', fontSize: 14, width: '100%' }}>
      {children}
    </Text>
  );
}
