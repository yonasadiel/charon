import React from 'react';
import { Field, reduxForm, InjectedFormProps, WrappedFieldProps } from 'redux-form';
import { Button, TextInput, BUTTON_TYPE_CONTAINED } from 'react-hephaestus';

import LoadingCircle from '../../loading/Circle';
import './LoginForm.scss';

export interface LoginFormData {
  username: string;
  password: string;
};

const renderUsernameField = (field: WrappedFieldProps) => (
  <TextInput {...field.input} helpText={field.meta.error} placeholder="username" type="username" />
);

const renderPasswordField = (field: WrappedFieldProps) => (
  <TextInput {...field.input} helpText={field.meta.error} placeholder="password" type="password" />
);

const LoginForm = (props: InjectedFormProps<LoginFormData>) => {
  const { error, handleSubmit, submitting } = props;
  return (
    <form className="login-form" onSubmit={handleSubmit}>
      <Field name="username" component={renderUsernameField} />
      <Field name="password" component={renderPasswordField} />

      <small className="error">{ error }</small>

      <div className="login-button-row">
        <Button buttonType={BUTTON_TYPE_CONTAINED} onClick={handleSubmit}>
          { !submitting ? 'LOGIN' : <LoadingCircle /> }
        </Button>
      </div>
    </form>
  );
};

export default reduxForm<LoginFormData>({ form: 'login' })(LoginForm);
