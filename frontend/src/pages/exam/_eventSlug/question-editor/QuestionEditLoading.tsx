import React from 'react';
import { Button } from 'react-hephaestus';

import './QuestionEditorPage.scss';


const QuestionEditLoading = () => {
  return (
    <div className="question">
      <div className="content"><span className="skeleton">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras elit enim, commodo vitae lectus pharetra, dictum semper enim. Ut eu erat a turpis aliquet luctus sed auctor dui. Aenean laoreet enim et iaculis porttitor. </span></div>
      <div className="editing-row">
        <Button className="edit-button skeleton"><i className="fas fa-edit"></i><strong>EDIT</strong></Button>
        <Button className="delete-button skeleton">
          <i className="fas fa-trash"></i><strong>DELETE</strong>
        </Button>
      </div>
    </div>
  );
};

export default QuestionEditLoading;