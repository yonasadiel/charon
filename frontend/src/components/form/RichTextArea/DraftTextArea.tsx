import { EditorState, convertToRaw, ContentState } from 'draft-js';
import draftToHtml from 'draftjs-to-html';
import htmlToDraft from 'html-to-draftjs';
import React from 'react';
import { Editor, EditorProps } from 'react-draft-wysiwyg';
import { WrappedFieldInputProps } from 'redux-form';

export interface DraftTextAreaProps extends WrappedFieldInputProps {
  wrapperClassName?: string;
  editorClassName?: string;
  toolbarClassName?: string;
}

const DraftTextArea = (props: DraftTextAreaProps & EditorProps) => {
  const { onChange, value, ...rest } = props;
  const contentState = ContentState.createFromBlockArray(htmlToDraft(value).contentBlocks);
  const initialEditorState = EditorState.createWithContent(contentState);
  const [editorState, setEditorState] = React.useState(initialEditorState);
  const handleEditorStateChange = (editorState: EditorState) => {
    const value = draftToHtml(convertToRaw(editorState.getCurrentContent()));
		setEditorState(editorState);
    onChange(value);
  };
  return (
    <Editor
      toolbar={{
        options: ['inline', 'fontSize', 'fontFamily', 'list', 'textAlign', 'image', 'remove', 'history'],
      }}
      editorState={editorState}
      onEditorStateChange={handleEditorStateChange}
      {...rest} />
  );

};

export default DraftTextArea;
