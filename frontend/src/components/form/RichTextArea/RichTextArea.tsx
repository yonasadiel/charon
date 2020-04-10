import * as React from 'react';

import { DraftTextAreaProps } from './DraftTextArea';
const LazyReactDraftEditor = React.lazy(async () => import('./DraftTextArea'));

export default (props: DraftTextAreaProps) => (
  <React.Suspense fallback={<textarea value={props.value} onChange={props.onChange} />}>
    <LazyReactDraftEditor {...props} />
  </React.Suspense>
);
