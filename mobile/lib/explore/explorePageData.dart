import 'package:flutter/material.dart';

enum SelectedCardListEnum {
    All,
    Workflows,
    Services,
}

class SelectedCardList {
  static SelectedCardListEnum cardType = SelectedCardListEnum.All;
  static List<Widget> cardList = [];
}