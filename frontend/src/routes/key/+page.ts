import type { PageLoad } from './$types';
import { WorkbooksMeta } from '../../../bindings/merger/services/workbook';

export const load: PageLoad = async () => {
    return {
        metas: await WorkbooksMeta(),
	};
};
