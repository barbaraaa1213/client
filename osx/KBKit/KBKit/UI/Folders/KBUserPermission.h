//
//  KBUserPermission.h
//  Keybase
//
//  Created by Gabriel on 4/30/15.
//  Copyright (c) 2015 Gabriel Handford. All rights reserved.
//

#import <Foundation/Foundation.h>

#import "KBRPC.h"

@interface KBUserPermission : NSObject

@property KBRUser *user;
@property NSString *permission;

@end
