//
//  KBDeviceSetupView.m
//  Keybase
//
//  Created by Gabriel on 2/18/15.
//  Copyright (c) 2015 Gabriel Handford. All rights reserved.
//

#import "KBDeviceSetupView.h"

#import "KBDeviceSignerView.h"
#import "KBRPC.h"
#import "KBDeviceSignerOption.h"

@interface KBDeviceSetupView ()
@end

@implementation KBDeviceSetupView

- (void)viewInit {
  [super viewInit];
  self.wantsLayer = YES;
  self.layer.backgroundColor = NSColor.whiteColor.CGColor;

  KBLabel *infoLabel = [[KBLabel alloc] init];
  [infoLabel setText:@"This is the first time you've logged into this computer. You need to setup and verify this installation of Keybase. Which method do you want to use?" font:[NSFont systemFontOfSize:14] color:[KBLookAndFeel textColor] alignment:NSCenterTextAlignment];
  [self addSubview:infoLabel];

  _deviceSignerView = [[KBDeviceSignerView alloc] init];
  [self addSubview:_deviceSignerView];

  YONSView *bottomView = [[YONSView alloc] init];
  _selectButton = [KBButton buttonWithText:@"Select" style:KBButtonStylePrimary];
  [bottomView addSubview:_selectButton];
  _cancelButton = [KBButton buttonWithText:@"Cancel" style:KBButtonStyleDefault];
  [bottomView addSubview:_cancelButton];
  [self addSubview:bottomView];

  YOSelf yself = self;
  bottomView.viewLayout = [YOLayout layoutWithLayoutBlock:^CGSize(id<YOLayout> layout, CGSize size) {
    CGFloat y = 0;
    y += [layout sizeToFitVerticalInFrame:CGRectMake(size.width - 280, 0, 130, 0) view:yself.selectButton].size.height;
    [layout sizeToFitVerticalInFrame:CGRectMake(size.width - 130, 0, 130, 0) view:yself.cancelButton];
    return CGSizeMake(size.width, y);
  }];

  self.viewLayout = [YOLayout layoutWithLayoutBlock:[self borderLayoutWithCenterView:_deviceSignerView topView:infoLabel bottomView:bottomView margin:UIEdgeInsetsMake(20, 20, 20, 20) padding:20]];
}

- (void)setDevices:(NSArray *)devices hasPGP:(BOOL)hasPGP {
  NSMutableArray *deviceSignerOptions = [NSMutableArray array];
  for (KBRDeviceDescription *device in devices) {
    [deviceSignerOptions addObject:[KBDeviceSignerOption deviceWithIdentifier:device.deviceID name:device.name]];
  }
  if (hasPGP) [deviceSignerOptions addObject:[KBDeviceSignerOption PGP]];

  [_deviceSignerView setDeviceSignerOptions:deviceSignerOptions];
}

- (YOLayoutBlock)borderLayoutWithCenterView:(id)centerView topView:(id)topView bottomView:(id)bottomView margin:(UIEdgeInsets)margin padding:(CGFloat)padding {
  return ^CGSize(id<YOLayout> layout, CGSize size) {

    CGSize sizeWithMargin = CGSizeMake(size.width - margin.left - margin.right, size.height - margin.top - margin.bottom);

    CGSize topSize = [topView sizeThatFits:sizeWithMargin];
    CGSize bottomSize = [bottomView sizeThatFits:sizeWithMargin];

    CGFloat centerHeight = sizeWithMargin.height - topSize.height - bottomSize.height - (padding * 2);

    CGFloat y = margin.top;
    y += [layout setFrame:CGRectMake(margin.left, y, topSize.width, topSize.height) view:topView].size.height + padding;

    y += [layout setFrame:CGRectMake(margin.left, y, sizeWithMargin.width, centerHeight) view:centerView].size.height + padding;

    y += [layout setFrame:CGRectMake(margin.left, y, bottomSize.width, bottomSize.height) view:bottomView].size.height;

    return size;
  };
}

@end
