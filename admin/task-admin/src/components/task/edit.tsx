import React from 'react'
import {BooleanInput, Edit, SelectInput, SimpleForm, TextInput, useInput} from 'react-admin'
import {WebhookTypes} from "./webhookType";

const TaskEdit = (props) => {
    return (
        <Edit title='Edit Task' {...props}>
            <SimpleForm>
                <TextInput disabled source='id'/>
                <TextInput source='name' required fullWidth/>
                <BooleanInput source='is_active'/>
                <TextInput source='schedule' required fullWidth/>
                <SelectInput choices={WebhookTypes} required source='webhook_type'/>
                <TextInput multiline fullWidth required source='webhook'/>
                <TextInput source='message' multiline fullWidth required/>
            </SimpleForm>
        </Edit>
    )
}

export default TaskEdit