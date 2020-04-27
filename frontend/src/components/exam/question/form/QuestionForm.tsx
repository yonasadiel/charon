import React from 'react';
import { Field, FieldArray, reduxForm, InjectedFormProps, WrappedFieldProps, WrappedFieldArrayProps, FieldArrayFieldsProps } from 'redux-form';
import { Button, TextInput, BUTTON_TYPE_CONTAINED } from 'react-hephaestus';

import { Question } from '../../../../modules/charon/exam/api';
import RichTextArea from '../../../form/RichTextArea/RichTextArea';
import LoadingCircle from '../../../loading/Circle';
import './QuestionForm.scss';

export interface QuestionFormData extends Partial<Question> {};

const renderContentField = (field: WrappedFieldProps) => {
  const { input } = field;
  return (
    <div className="field-row">
      <p><strong>Pertanyaan</strong></p>
      <RichTextArea {...input} editorClassName="editor-area" />
      <small>{field.meta.error}</small>
    </div>
  );
}

const renderChoiceField = (field: WrappedFieldProps) => (
  <TextInput {...field.input} errorText={field.meta.error} type="text" />
);

const renderChoice = (member: string, index: number, fields: FieldArrayFieldsProps<string>) => (
  <div className="field-row choice" key={index}>
    <div className="label"><strong>Jawaban #{index + 1}</strong></div>
    <Field name={`${member}`} component={renderChoiceField} />
    <Button onClick={() => fields.remove(index)} className="delete-button"><i className="fas fa-trash"></i></Button>
  </div>
);

const renderChoices = (props: WrappedFieldArrayProps<string>) => {
  const { fields } = props;
  const addFourMoreChoices = () => {
    fields.push('');
    fields.push('');
    fields.push('');
    fields.push('');
  }
  return (
    <div>
      {fields.map(renderChoice)}
      <div className="field-row add-choice">
        <Button buttonType="outlined" onClick={addFourMoreChoices}><i className="fas fa-plus"></i>Tambah pilihan jawaban</Button>
      </div>
    </div>
  );
}

const QuestionForm = (props: InjectedFormProps<QuestionFormData>) => {
  const { error, handleSubmit, submitting, initialValues } = props;
  return (
    <form className="question-form" onSubmit={handleSubmit}>
      <h2 className="title">{!!initialValues ? 'Ubah Pertanyaan' : 'Tambah Pertanyaan'}</h2>
      <Field name="content" component={renderContentField} />
      <FieldArray name="choices" component={renderChoices} />

      <small className="error">{ error }</small>

      <div className="submit-button-row">
        <Button buttonType={BUTTON_TYPE_CONTAINED} type="submit">
          { !submitting ? (<strong>BUAT</strong>) : <LoadingCircle /> }
        </Button>
      </div>
    </form>
  );
};

export default reduxForm<QuestionFormData>({ form: 'question' })(QuestionForm);
