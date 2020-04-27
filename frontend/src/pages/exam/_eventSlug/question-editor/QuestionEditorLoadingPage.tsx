import React from 'react';
import { Card } from 'react-hephaestus';
import QuestionEditLoading from './QuestionEditLoading';

const QuestionEditorLoadingPage = () => (
  <Card className="question-editor">
    <QuestionEditLoading />
    <QuestionEditLoading />
    <QuestionEditLoading />
  </Card>
);

export default QuestionEditorLoadingPage;
