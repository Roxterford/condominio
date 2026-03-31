import { Checkbox } from "@/components/ui/checkbox";
import { DatePickerInput } from "@/components/ui/date-picker-input";
import { Input } from "@/components/ui/input";
import { InputGroupInput } from "@/components/ui/input-group";
import { Select } from "@/components/ui/select";
import { createFormHook, createFormHookContexts } from "@tanstack/react-form";

const { fieldContext, formContext } = createFormHookContexts();

// Allow us to bind components to the form to keep type safety but reduce production boilerplate
// Define this once to have a generator of consistent form instances throughout your app
export const { useAppForm, withForm } = createFormHook({
  fieldComponents: {
    Input,
    Select,
    Checkbox,
    DatePickerInput,
    InputGroupInput,
  },
  formComponents: {},
  fieldContext,
  formContext,
});
