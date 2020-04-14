/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package br.com.zup.charles.domain

enum class DeploymentStatusEnum {
    NOT_DEPLOYED,
    DEPLOYING,
    DEPLOYED,
    DEPLOY_FAILED,
    UNDEPLOYING
}