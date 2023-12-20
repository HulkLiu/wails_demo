import { type PropType } from 'vue';
import { type RowKey } from '../interface';
declare const _default: import("vue").DefineComponent<{
    rowKey: {
        type: PropType<RowKey>;
        required: true;
    };
    disabled: {
        type: BooleanConstructor;
        required: true;
    };
    onUpdateChecked: {
        type: PropType<(checked: boolean, e: MouseEvent | KeyboardEvent) => void>;
        required: true;
    };
}, () => JSX.Element, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
    rowKey: {
        type: PropType<RowKey>;
        required: true;
    };
    disabled: {
        type: BooleanConstructor;
        required: true;
    };
    onUpdateChecked: {
        type: PropType<(checked: boolean, e: MouseEvent | KeyboardEvent) => void>;
        required: true;
    };
}>>, {}, {}>;
export default _default;
