import { AuthService } from "./auth/v1/auth_connect";
import * as AuthPb from "./auth/v1/auth_pb";

export default {
    Services: { AuthService },
    Pb: { AuthPb },
};
