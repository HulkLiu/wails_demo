import { type Ref } from 'vue';
import { type TreeNode } from 'treemate';
import { type Key, type TmNode, type TreeOption } from './interface';
export declare function useKeyboard({ props, fNodesRef, mergedExpandedKeysRef, mergedSelectedKeysRef, handleSelect, handleSwitcherClick }: {
    props: {
        keyboard: boolean;
    };
    fNodesRef: Ref<Array<TreeNode<TreeOption>>>;
    mergedExpandedKeysRef: Ref<Key[]>;
    mergedSelectedKeysRef: Ref<Key[]>;
    handleSelect: (node: TmNode) => void;
    handleSwitcherClick: (node: TmNode) => void;
}): {
    pendingNodeKeyRef: Ref<null | Key>;
    handleKeydown: (e: KeyboardEvent) => void;
};
