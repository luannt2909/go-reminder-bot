import React from 'react'
import { Create, SimpleForm, TextInput, DateInput, SelectInput } from 'react-admin'
import { RichTextInput } from 'ra-input-rich-text';
import {WebhookTypes} from './webhookType'

const TaskCreate = (props) => {
    return (
        <Create title='Create a Task' {...props}>
            <SimpleForm>
                <TextInput source='name' />
                <TextInput multiline source='schedule' />
                <SelectInput choices={WebhookTypes} source='webhook_type' />
                <TextInput multiline source='webhook' />
                <RichTextInput multiline source='message' />
            </SimpleForm>
        </Create>
    )
}

export default TaskCreate