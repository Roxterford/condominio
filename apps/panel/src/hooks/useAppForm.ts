import { Checkbox } from "@/components/ui/checkbox";
import { DatePickerInput } from "@/components/ui/date-picker-input";
import { Input } from "@/components/ui/input";
import {
  createFormHook,
  createFormHookContexts,
} from "@tanstack/react-form-nextjs";

const { fieldContext, formContext } = createFormHookContexts();

const { useAppForm } = createFormHook({
  fieldComponents: {
    DatePickerInput,
    Input,
    Checkbox,
  },
  formComponents: {},
  fieldContext,
  formContext,
});

export { useAppForm };
